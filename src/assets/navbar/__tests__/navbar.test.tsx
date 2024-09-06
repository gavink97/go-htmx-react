/**
 * @jest-environment jsdom
 */

// import { render, fireEvent, screen } from 'test-utils';
import { fireEvent, render, screen, waitFor} from '@testing-library/react'
import userEvent from '@testing-library/user-event';
import { NavBar } from '../navbar';
import '@testing-library/jest-dom';

test('renders NavBar and interacts with the dropdown menu', async () => {
    render(<NavBar />);

    const menuTrigger = screen.getByRole('button', { name: /menu/i });
    expect(menuTrigger).toBeInTheDocument();

    const dropdownContent = screen.queryByText(/my account/i);
    expect(dropdownContent).not.toBeInTheDocument();

    await userEvent.click(menuTrigger);

    await waitFor(() => {
        expect(screen.getByRole('menu')).toBeInTheDocument();
    });
    // screen.debug();
    const menuItems = ['My Account', 'Profile', 'Billing', 'Team', 'Subscription'];
    menuItems.forEach((item) => {
        expect(screen.queryByText(item)).toBeInTheDocument();
    });

    expect(menuTrigger).toHaveAttribute('aria-expanded', 'true');

    fireEvent.click(menuTrigger);
    expect(dropdownContent).not.toBeInTheDocument();
});
